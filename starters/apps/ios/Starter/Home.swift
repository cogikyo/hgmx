import SwiftUI

struct Home: View {
    let session: Session

    var body: some View {
        if let account = session.account {
            Centered {
                Text("Hi, \(account.displayName)").font(.title2)
                NavigationLink("Account", value: Route.account)
            }
        }
    }
}
