package server

// type RegistryPlugin struct {
// 	acts *hashset.Set
// }

// func NewRegistryPlugin() *RegistryPlugin {
// 	return &RegistryPlugin{}
// }

// func (plugin *RegistryPlugin) OnStart(ctx actor.Context) {
// 	switch ctx.Actor().(type) {
// 	case server.Server:
// 		log.Info("Registry plugin: server started")
// 	case server.Record:
// 		log.Info("Registry plugin: registry actor")

// 		plugin.acts.Add(ctx.Actor())
// 	}
// }

// func (plugin *RegistryPlugin) OnOtherMessage(ctx actor.Context, msg interface{}) {
// 	switch msg.(type) {
// 	case *actor.Stopped:
// 		if _, ok := ctx.Actor().(server.Record); ok {
// 			plugin.acts.Remove(ctx.Actor())
// 		}
// 	}
// }
