import API
import SwiftUI

@MainActor
@Observable
final class NotesModel {
    let work = Work()
    var draft = ""
    private(set) var notes: Fixture_Api_V1_Notes?
    private let session: Session
    private let client: Fixture_Api_V1_NoteServiceClient

    init(session: Session) {
        self.session = session
        client = Fixture_Api_V1_NoteServiceClient(client: session.api.client)
    }

    func load() async {
        await work.run {
            notes = try await session.call { await self.client.listNotes(request: .init(), headers: $0) }.notes
        }
    }

    func add() async {
        await work.run {
            var request = Fixture_Api_V1_AddNoteRequest()
            request.text = draft.trimmingCharacters(in: .whitespaces)
            notes = try await session.call { await self.client.addNote(request: request, headers: $0) }.notes
            draft = ""
        }
    }
}

struct Home: View {
    @State private var model: NotesModel

    init(session: Session) {
        _model = State(initialValue: NotesModel(session: session))
    }

    var body: some View {
        @Bindable var model = model
        List {
            ProblemText(problem: model.work.problem)
            if let notes = model.notes {
                if notes.canAdd {
                    HStack {
                        TextField("New note", text: $model.draft)
                        Button("Add") { Task { await model.add() } }
                            .disabled(model.draft.trimmingCharacters(in: .whitespaces).isEmpty)
                    }
                }
                if notes.notes.isEmpty {
                    Text("No notes yet.")
                }
                ForEach(notes.notes, id: \.id) { Text($0.text) }
            } else if model.work.busy {
                ProgressView()
            }
        }
        .disabled(model.work.busy)
        .navigationTitle("Notes")
        .toolbar { NavigationLink("Account", value: Route.account) }
        .task { await model.load() }
    }
}
