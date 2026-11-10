import Observation

@MainActor
@Observable
final class Work {
    private(set) var busy = false
    private(set) var problem: Problem?

    func run(_ action: () async throws -> Void) async {
        guard !busy else { return }
        busy = true
        problem = nil
        defer { busy = false }
        do {
            try await action()
        } catch let p as Problem {
            problem = p
        } catch {}
    }
}
