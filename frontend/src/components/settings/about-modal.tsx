import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { GetVersion } from "@wails/go/main/App";
import { useEffect, useState } from "react";

interface AboutModalProps {
  open: boolean;
  onClose: () => void;
}

export function AboutModal({ open, onClose }: AboutModalProps) {
  const [version, setVersion] = useState("");

  useEffect(() => {
    if (open) {
      GetVersion().then((v: string) => setVersion(v)).catch(() => setVersion(""));
    }
  }, [open]);

  return (
    <Dialog open={open} onOpenChange={(v) => { if (!v) onClose(); }}>
      <DialogContent className="max-w-sm border-border bg-surface text-text-primary">
        <DialogHeader>
          <DialogTitle className="font-mono text-lg">Ignite</DialogTitle>
        </DialogHeader>
        <div className="flex flex-col items-center space-y-3 text-sm text-text-secondary">
          <p className="font-mono text-base text-text-primary">
            Provisioning with a heartbeat
          </p>
          <p className="font-mono text-sm text-accent">
            v{version}
          </p>
          <p className="w-full max-w-xs text-left">
            A desktop GUI for provisioning new software
            documentation quality. Conducts AI-guided interviews and generates
            project specs, agent guides, implementation plans, and READMEs.
          </p>
          <div className="w-full max-w-xs space-y-1 pt-2 text-left">
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
