import API
import SwiftUI

@MainActor
@Observable
final class SessionsModel {
    let work = Work()
    private(set) var devices: [Hgmx_Auth_V1_Session]?
    private let session: Session

    init(session: Session) {
        self.session = session
    }

    func load() async {
        await work.run { try await list() }
    }

    func revoke(_ device: Hgmx_Auth_V1_Session) async {
        await work.run {
            var request = Hgmx_Auth_V1_RevokeSessionRequest()
            request.sessionID = device.id
            _ = try await session.call { await self.session.api.auth.revokeSession(request: request, headers: $0) }
            try await list()
        }
    }

    func revokeOthers() async {
        await work.run {
            _ = try await session.call { await self.session.api.auth.revokeOtherSessions(request: .init(), headers: $0) }
            try await list()
        }
    }

    private func list() async throws {
        devices = try await session.call { await self.session.api.auth.listSessions(request: .init(), headers: $0) }.sessions
    }
}

struct Sessions: View {
    @State private var model: SessionsModel

    init(session: Session) {
        _model = State(initialValue: SessionsModel(session: session))
    }

    var body: some View {
        List {
            ProblemText(problem: model.work.problem)
            if let devices = model.devices {
                ForEach(devices, id: \.id) { device in
                    HStack {
                        VStack(alignment: .leading) {
                            Text(device.label.isEmpty ? kind(device.kind) : device.label)
                            Text("\(kind(device.kind)) · last seen \(date(device.seenAt))")
                                .font(.footnote)
                                .foregroundStyle(.secondary)
                        }
                        Spacer()
                        if device.current {
                            Text("This device").foregroundStyle(.secondary)
                        } else {
                            Button("Sign out") { Task { await model.revoke(device) } }
                        }
                    }
                }
                if devices.contains(where: { !$0.current }) {
                    Button("Sign out everywhere else") { Task { await model.revokeOthers() } }
                }
            } else if model.work.busy {
                ProgressView()
            }
        }
        .disabled(model.work.busy)
        .navigationTitle("Signed-in devices")
        .task { await model.load() }
    }

    private func kind(_ kind: Hgmx_Auth_V1_ClientKind) -> String {
        switch kind {
        case .browser: "Browser"
        case .ios: "iPhone or iPad"
        case .android: "Android"
        default: "Other device"
        }
    }
}
