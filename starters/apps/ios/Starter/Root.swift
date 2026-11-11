import SwiftUI

enum Route: Hashable {
    case account
    case sessions
    case reauth
}

struct Root: View {
    let session: Session

    var body: some View {
        switch session.gate {
        case .starting:
            ProgressView()
        case .unreachable(let problem):
            Centered {
                Text(problem.message)
                Button("Try again") { Task { await session.start() } }
            }
        case .update(let store, let message):
            Centered {
                Text("Update required").font(.title2)
                Text(message)
                if let warning = session.warning { Text(warning) }
                if let url = URL(string: store), url.scheme == "https" {
                    Link("Open the store", destination: url)
                }
            }
        case .signedOut(let notice):
            SignIn(session: session, notice: notice, warning: session.warning)
        case .signedIn(let account):
            SignedIn(session: session).id(account.id)
        }
    }
}

private struct SignedIn: View {
    let session: Session
    @State private var path: [Route] = []

    var body: some View {
        NavigationStack(path: $path) {
            Home(session: session)
                .navigationDestination(for: Route.self) { route in
                    switch route {
                    case .account:
                        Account(session: session)
                    case .sessions:
                        Sessions(session: session)
                    case .reauth:
                        SignIn(session: session, notice: "Sign in again to continue.") { path.removeLast() }
                    }
                }
        }
    }
}

struct Centered<Content: View>: View {
    @ViewBuilder let content: Content

    var body: some View {
        VStack(spacing: 16) {
            content
        }
        .multilineTextAlignment(.center)
        .padding(24)
    }
}

struct ProblemText: View {
    let problem: Problem?

    var body: some View {
        if let problem {
            Text(problem.message).foregroundStyle(.red)
            if problem.reauth {
                NavigationLink("Sign in again", value: Route.reauth)
            }
        }
    }
}
