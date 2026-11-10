import API
import SwiftProtobuf
import SwiftUI

@MainActor
@Observable
final class AccountModel {
    let work = Work()
    private(set) var notice: String?
    private let session: Session

    init(session: Session) {
        self.session = session
    }

    func username(_ name: String) async {
        await act {
            var request = Hgmx_Auth_V1_SetUsernameRequest()
            request.username = name
            try await session.update({ await self.session.api.auth.setUsername(request: request, headers: $0) }) { $0.account }
        }
    }

    func email(_ email: String) async {
        await act {
            var request = Hgmx_Auth_V1_ChangeEmailRequest()
            request.email = email
            _ = try await session.call { await self.session.api.auth.changeEmail(request: request, headers: $0) }
            notice = "Open the link we sent to \(email) to finish the change."
        }
    }

    func addPasskey() async {
        await act {
            let begin = try await session.call { await self.session.api.auth.beginPasskeyAdd(request: .init(), headers: $0) }
            var request = Hgmx_Auth_V1_FinishPasskeyAddRequest()
            request.ceremony = begin.options.ceremony
            request.credentialJson = try await Credentials.createPasskey(begin.options.publicKeyOptionsJson)
            try await session.update({ await self.session.api.auth.finishPasskeyAdd(request: request, headers: $0) }) { $0.account }
        }
    }

    func rename(_ passkey: Hgmx_Auth_V1_Passkey, to name: String) async {
        await act {
            var request = Hgmx_Auth_V1_RenamePasskeyRequest()
            request.passkeyID = passkey.id
            request.name = name
            try await session.update({ await self.session.api.auth.renamePasskey(request: request, headers: $0) }) { $0.account }
        }
    }

    func remove(_ passkey: Hgmx_Auth_V1_Passkey) async {
        await act {
            var request = Hgmx_Auth_V1_RemovePasskeyRequest()
            request.passkeyID = passkey.id
            try await session.update({ await self.session.api.auth.removePasskey(request: request, headers: $0) }) { $0.account }
        }
    }

    func unlinkGoogle() async {
        await act {
            try await session.update({ await self.session.api.auth.unlinkGoogle(request: .init(), headers: $0) }) { $0.account }
        }
    }

    func delete() async {
        await act { try await session.deleteAccount() }
    }

    private func act(_ action: () async throws -> Void) async {
        notice = nil
        await work.run(action)
    }
}

struct Account: View {
    let session: Session
    @State private var model: AccountModel
    @State private var username = ""
    @State private var email = ""
    @State private var renaming: Hgmx_Auth_V1_Passkey?
    @State private var name = ""
    @State private var deleting = false

    init(session: Session) {
        self.session = session
        _model = State(initialValue: AccountModel(session: session))
    }

    var body: some View {
        if let account = session.account {
            form(account)
                .navigationTitle("Account")
                .onAppear { username = account.username }
        }
    }

    private func form(_ account: Hgmx_Auth_V1_Account) -> some View {
        let caps = account.capabilities
        let google = account.identities.first { $0.provider == .google }
        return Form {
            Section {
                ProblemText(problem: model.work.problem)
                if let notice = model.notice {
                    Text(notice)
                }
                Text(account.displayName).font(.title2)
                if account.hasEmail {
                    Text(account.email)
                }
            }
            Section {
                TextField("Username", text: $username)
                    .textInputAutocapitalization(.never)
                    .autocorrectionDisabled()
                Button("Save username") { Task { await model.username(username.trimmingCharacters(in: .whitespaces)) } }
                    .disabled(username.trimmingCharacters(in: .whitespaces) == account.username)
            }
            if caps.changeEmail {
                Section {
                    TextField("New email", text: $email)
                        .keyboardType(.emailAddress)
                        .textInputAutocapitalization(.never)
                        .autocorrectionDisabled()
                    Button("Change email") { Task { await model.email(email.trimmingCharacters(in: .whitespaces)) } }
                        .disabled(email.trimmingCharacters(in: .whitespaces).isEmpty)
                }
            }
            Section("Passkeys") {
                ForEach(account.passkeys, id: \.id) { passkey in
                    VStack(alignment: .leading) {
                        Text(passkey.name)
                        Text(detail(passkey)).font(.footnote).foregroundStyle(.secondary)
                    }
                    .swipeActions {
                        if passkey.removable {
                            Button("Remove", role: .destructive) { Task { await model.remove(passkey) } }
                        }
                        Button("Rename") {
                            name = passkey.name
                            renaming = passkey
                        }
                    }
                }
                Button("Add a passkey") { Task { await model.addPasskey() } }
            }
            if caps.google, let google {
                Section {
                    if google.removable {
                        Button("Unlink Google") { Task { await model.unlinkGoogle() } }
                    } else {
                        Text("Google is your only way to sign in.")
                    }
                }
            }
            ForEach(Array(account.identities.filter { $0.provider != .google }.indices), id: \.self) { _ in
                Text("Another sign-in method")
            }
            Section {
                NavigationLink("Signed-in devices", value: Route.sessions)
                Button("Sign out") { session.signOut() }
                if caps.deleteAccount {
                    Button("Delete account", role: .destructive) { deleting = true }
                }
            }
        }
        .disabled(model.work.busy)
        .alert("Rename passkey", isPresented: Binding { renaming != nil } set: { if !$0 { renaming = nil } }) {
            TextField("Name", text: $name)
            Button("Save") {
                if let passkey = renaming {
                    Task { await model.rename(passkey, to: name.trimmingCharacters(in: .whitespaces)) }
                }
            }
            .disabled(name.trimmingCharacters(in: .whitespaces).isEmpty)
            Button("Cancel", role: .cancel) {}
        }
        .confirmationDialog("Delete your account?", isPresented: $deleting, titleVisibility: .visible) {
            Button("Delete", role: .destructive) { Task { await model.delete() } }
        } message: {
            Text("This removes your account and signs out every device. It can't be undone.")
        }
    }

    private func detail(_ passkey: Hgmx_Auth_V1_Passkey) -> String {
        var parts = ["Added \(date(passkey.createdAt))"]
        if passkey.hasUsedAt {
            parts.append("last used \(date(passkey.usedAt))")
        }
        if !passkey.usable {
            parts.append("not usable here")
        }
        return parts.joined(separator: " · ")
    }
}

func date(_ t: Google_Protobuf_Timestamp) -> String {
    t.date.formatted(date: .abbreviated, time: .omitted)
}
