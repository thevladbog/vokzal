import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { FluentProvider, webLightTheme } from '@fluentui/react-components';
import { PublicBoardPage } from '@/pages/PublicBoardPage';
import { PlatformBoardPage } from '@/pages/PlatformBoardPage';

function App() {
  return (
    <FluentProvider theme={webLightTheme}>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<PublicBoardPage />} />
          <Route path="/platform" element={<PlatformBoardPage />} />
        </Routes>
      </BrowserRouter>
    </FluentProvider>
  );
}

export default App;
