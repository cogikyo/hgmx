import SwiftUI

@main
struct StarterApp: App {
    @Environment(\.scenePhase) private var phase
    @State private var session: Session
    @State private var away = false

    init() {
        let info = Bundle.main
        let api = Api(
            origin: info.object(forInfoDictionaryKey: "HGMXAPIOrigin") as? String ?? "",
            appID: info.object(forInfoDictionaryKey: "HGMXAppID") as? String ?? "",
            build: info.object(forInfoDictionaryKey: "CFBundleVersion") as? String ?? ""
        )
        _session = State(initialValue: Session(api: api, tokens: TokenStore(origin: api.origin)))
    }

    var body: some Scene {
        WindowGroup {
            Root(session: session)
                .task { await session.start() }
        }
        .onChange(of: phase) { _, phase in
            switch phase {
            case .background:
                away = true
            case .active where away:
                away = false
                Task { await session.recheck() }
            default:
                break
            }
        }
    }
}
