import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

interface AboutModalProps {
  open: boolean;
  onClose: () => void;
}

export function AboutModal({ open, onClose }: AboutModalProps) {
  return (
    <Dialog open={open} onOpenChange={(v) => { if (!v) onClose(); }}>
      <DialogContent className="max-w-sm border-border bg-surface text-text-primary">
        <DialogHeader>
          <DialogTitle className="font-mono text-lg">Ignite</DialogTitle>
        </DialogHeader>
        <div className="space-y-3 text-sm text-text-secondary">
          <p className="font-mono text-base text-text-primary">
            Provisioning with a heartbeat
          </p>
          <p>
            A desktop GUI for provisioning new software projects with Kvasir-level
            documentation quality. Conducts AI-guided interviews and generates
            project specs, agent guides, implementation plans, and READMEs.
          </p>
          <div className="space-y-1 pt-2">
            <p>
              <span className="text-text-primary">Author:</span> Morten Johansen
            </p>
            <p>
              <span className="text-text-primary">Web:</span>{" "}
              <a
                href="https://johansen.foo"
                target="_blank"
                rel="noopener noreferrer"
                className="text-accent hover:underline"
              >
                johansen.foo
              </a>
            </p>
            <p>
              <span className="text-text-primary">Repo:</span>{" "}
              <a
                href="https://github.com/mojoaar/ignite"
                target="_blank"
                rel="noopener noreferrer"
                className="text-accent hover:underline"
              >
                github.com/mojoaar/ignite
              </a>
            </p>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  );
}
