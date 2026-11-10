import SwiftUI

@main
struct StarterApp: App {
    @State private var session: Session

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
    }
}
