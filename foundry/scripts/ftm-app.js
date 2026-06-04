export class FtmApp extends foundry.applications.api.HandlebarsApplicationMixin(
  foundry.applications.api.ApplicationV2,
) {
  static DEFAULT_OPTIONS = {
    id: 'ftm-tunnel-app',
    title: 'FTM Tunnel Manager',
    width: 960,
    height: 720,
    resizable: true,
    template: 'modules/ftm-tunnel/templates/app.hbs',
  };

  constructor(manager) {
    super();
    this.manager = manager;
  }

  async _prepareContext() {
    return { apiBase: this.manager?.apiBase ?? 'http://localhost:40500' };
  }
}
