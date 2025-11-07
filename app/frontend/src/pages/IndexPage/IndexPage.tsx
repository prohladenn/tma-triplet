import {
  Button,
  Cell,
  Input,
  List,
  Modal,
  Placeholder,
  Section,
  Text,
} from "@telegram-apps/telegram-ui";
import type { FC } from "react";
import { useEffect, useRef, useState } from "react";

import { Link } from "@/components/Link/Link.tsx";
import { Page } from "@/components/Page.tsx";
import { useNotes } from "@/hooks/useNotes";

export const IndexPage: FC = () => {
  const {
    notes,
    isLoading,
    error,
    isBackendAvailable,
    addNote,
    deleteNote,
    deleteAllNotes,
  } = useNotes();
  const [newNoteText, setNewNoteText] = useState("");
  const [showBackendWarning, setShowBackendWarning] = useState(false);
  const [showBackendSuccess, setShowBackendSuccess] = useState(false);
  const prevBackendStatus = useRef<boolean | null>(null);

  // Show modal when backend status changes
  useEffect(() => {
    console.log("[IndexPage] Backend status check:", {
      isLoading,
      isBackendAvailable,
      prevStatus: prevBackendStatus.current,
    });

    if (!isLoading) {
      // Only show modals if status has changed (not on first load if backend is working)
      if (prevBackendStatus.current !== null) {
        if (prevBackendStatus.current !== isBackendAvailable) {
          console.log(
            "[IndexPage] Backend status changed from",
            prevBackendStatus.current,
            "to",
            isBackendAvailable
          );
          // Status changed - show appropriate modal
          if (!isBackendAvailable) {
            console.log("[IndexPage] Showing backend unavailable warning");
            setShowBackendWarning(true);
          } else {
            console.log("[IndexPage] Showing backend recovered success");
            setShowBackendSuccess(true);
          }
        }
      } else if (!isBackendAvailable) {
        // First time and backend is unavailable - show warning
        console.log(
          "[IndexPage] First load: backend unavailable, showing warning"
        );
        setShowBackendWarning(true);
      } else {
        // First time and backend is available - no modal (everything is working)
        console.log(
          "[IndexPage] First load: backend available, no modal needed"
        );
      }

      // Update previous status
      prevBackendStatus.current = isBackendAvailable;
    }
  }, [isBackendAvailable, isLoading]);

  const handleAddNote = async () => {
    const trimmedText = newNoteText.trim();
    if (trimmedText) {
      await addNote(trimmedText);
      setNewNoteText("");
    }
  };

  const formatDate = (timestamp: number) => {
    const date = new Date(timestamp);
    return date.toLocaleString();
  };

  return (
    <Page back={false}>
      <Modal
        open={showBackendWarning}
        onOpenChange={setShowBackendWarning}
        header={<Modal.Header>Backend Unavailable</Modal.Header>}
      >
        <Placeholder
          header="Working in Offline Mode"
          description="The backend server is currently unavailable. Your notes will be saved locally on this device only. They will sync to the server once it becomes available."
          action={
            <Button
              size="m"
              stretched
              onClick={() => setShowBackendWarning(false)}
            >
              Continue
            </Button>
          }
        >
          <div style={{ fontSize: "48px" }}>⚠️</div>
        </Placeholder>
      </Modal>

      <Modal
        open={showBackendSuccess}
        onOpenChange={setShowBackendSuccess}
        header={<Modal.Header>Backend Connected</Modal.Header>}
      >
        <Placeholder
          header="All Systems Operational"
          description="Successfully connected to the backend server. Your notes are being saved securely and will sync across all your devices."
          action={
            <Button
              size="m"
              stretched
              onClick={() => setShowBackendSuccess(false)}
            >
              Great!
            </Button>
          }
        >
          <div style={{ fontSize: "48px" }}>✅</div>
        </Placeholder>
      </Modal>

      <List>
        <Section
          header="Application Launch Data"
          footer="These pages help developer to learn more about current launch information"
        >
          <Link to="/init-data">
            <Cell subtitle="User data, chat information, technical data">
              Init Data
            </Cell>
          </Link>
          <Link to="/launch-params">
            <Cell subtitle="Platform identifier, Mini Apps version, etc.">
              Launch Parameters
            </Cell>
          </Link>
          <Link to="/theme-params">
            <Cell subtitle="Telegram application palette information">
              Theme Parameters
            </Cell>
          </Link>
        </Section>

        <Section
          header="Notes"
          footer={error || "Add notes that are saved to the backend"}
        >
          <Input
            header="New Note"
            placeholder="Enter your note..."
            value={newNoteText}
            onChange={(e) => setNewNoteText(e.target.value)}
            disabled={isLoading}
            after={
              <Button
                onClick={handleAddNote}
                disabled={!newNoteText.trim() || isLoading}
                loading={isLoading}
              >
                Add
              </Button>
            }
          />
          {notes.length > 0 && (
            <Cell>
              <Button
                size="s"
                mode="plain"
                stretched
                onClick={deleteAllNotes}
                disabled={isLoading}
              >
                Delete All Notes
              </Button>
            </Cell>
          )}
        </Section>

        {notes.length > 0 && (
          <Section header={`My Notes (${notes.length})`}>
            {notes.map((note) => (
              <Cell
                key={note.id}
                subtitle={formatDate(note.timestamp)}
                multiline
                after={
                  <Button
                    size="s"
                    mode="plain"
                    onClick={() => deleteNote(note.id)}
                    disabled={isLoading}
                  >
                    Delete
                  </Button>
                }
              >
                <Text style={{ wordBreak: "break-word" }}>{note.text}</Text>
              </Cell>
            ))}
          </Section>
        )}
      </List>
    </Page>
  );
};
