/** @type {import('vitepress').UserConfig} */
export default {
  title: 'Вокзал.ТЕХ Documentation',
  description: 'Система автоматизации автовокзала',
  base: '/vokzal/',
  // Ссылки на файлы вне docs/ (QUICKSTART, CONTRIBUTING, README в корне) и на каталоги
  // без index считаются «мёртвыми» при сборке, но ведут на существующие ресурсы в репозитории.
  ignoreDeadLinks: true,
  themeConfig: {
    nav: [
      { text: 'Главная', link: '/' },
      { text: 'API', link: '/api/' },
      { text: 'Руководства', link: '/user-guides/' },
    ],
    sidebar: 'auto',
  },
}
