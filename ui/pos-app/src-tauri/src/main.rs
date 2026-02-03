// Prevents additional console window on Windows in release, DO NOT REMOVE!!
#![cfg_attr(not(debug_assertions), windows_subsystem = "windows")]

use serde::{Deserialize, Serialize};
use tauri::Manager;

#[derive(Debug, Serialize, Deserialize)]
struct ApiError {
    error: String,
}

#[derive(Debug, Serialize, Deserialize)]
struct TicketSaleRequest {
    trip_id: String,
    passenger_fio: Option<String>,
    passenger_phone: Option<String>,
    seat_id: Option<String>,
}

#[derive(Debug, Serialize, Deserialize)]
struct Ticket {
    id: String,
    trip_id: String,
    price: f64,
    qr_code: String,
    bar_code: String,
}

#[derive(Debug, Serialize, Deserialize)]
struct Receipt {
    success: bool,
    fiscal_sign: String,
    ofd_url: String,
}

// Команда для продажи билета
#[tauri::command]
async fn sell_ticket(
    api_url: String,
    token: String,
    request: TicketSaleRequest,
) -> Result<Ticket, String> {
    let client = reqwest::Client::new();
    
    let response = client
        .post(format!("{}/tickets/sell", api_url))
        .header("Authorization", format!("Bearer {}", token))
        .json(&request)
        .send()
        .await
        .map_err(|e| format!("Network error: {}", e))?;

    if response.status().is_success() {
        let ticket: Ticket = response
            .json()
            .await
            .map_err(|e| format!("Parse error: {}", e))?;
        Ok(ticket)
    } else {
        let error: ApiError = response
            .json()
            .await
            .unwrap_or(ApiError {
                error: "Unknown error".to_string(),
            });
        Err(error.error)
    }
}

// Команда для возврата билета
#[tauri::command]
async fn return_ticket(
    api_url: String,
    token: String,
    ticket_id: String,
) -> Result<Ticket, String> {
    let client = reqwest::Client::new();
    
    let response = client
        .post(format!("{}/tickets/{}/return", api_url, ticket_id))
        .header("Authorization", format!("Bearer {}", token))
        .send()
        .await
        .map_err(|e| format!("Network error: {}", e))?;

    if response.status().is_success() {
        let ticket: Ticket = response
            .json()
            .await
            .map_err(|e| format!("Parse error: {}", e))?;
        Ok(ticket)
    } else {
        let error: ApiError = response
            .json()
            .await
            .unwrap_or(ApiError {
                error: "Unknown error".to_string(),
            });
        Err(error.error)
    }
}

// Команда для печати билета через локальный агент
#[tauri::command]
async fn print_ticket(
    agent_url: String,
    ticket_data: serde_json::Value,
) -> Result<bool, String> {
    let client = reqwest::Client::new();
    
    let response = client
        .post(format!("{}/printer/ticket", agent_url))
        .json(&ticket_data)
        .send()
        .await
        .map_err(|e| format!("Printer error: {}", e))?;

    Ok(response.status().is_success())
}

// Команда для печати чека через локальный агент
#[tauri::command]
async fn print_receipt(
    agent_url: String,
    receipt_data: serde_json::Value,
) -> Result<Receipt, String> {
    let client = reqwest::Client::new();
    
    let response = client
        .post(format!("{}/kkt/receipt", agent_url))
        .json(&receipt_data)
        .send()
        .await
        .map_err(|e| format!("KKT error: {}", e))?;

    if response.status().is_success() {
        let receipt: Receipt = response
            .json()
            .await
            .map_err(|e| format!("Parse error: {}", e))?;
        Ok(receipt)
    } else {
        Err("Failed to print receipt".to_string())
    }
}

// Команда для открытия окна экрана покупателя
#[tauri::command]
async fn open_customer_display(app_handle: tauri::AppHandle) -> Result<(), String> {
    tauri::WindowBuilder::new(
        &app_handle,
        "customer_display",
        tauri::WindowUrl::App("customer-display.html".into()),
    )
    .title("Экран покупателя")
    .inner_size(800.0, 600.0)
    .build()
    .map_err(|e| format!("Failed to open window: {}", e))?;

    Ok(())
}

fn main() {
    tauri::Builder::default()
        .invoke_handler(tauri::generate_handler![
            sell_ticket,
            return_ticket,
            print_ticket,
            print_receipt,
            open_customer_display
        ])
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
