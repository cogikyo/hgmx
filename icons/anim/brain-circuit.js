import { animate, stagger, τmin } from "../motion.js";

export default function (el) {
  const traces = ["M12 8h8", "M16 8V5", "M12 13h4", "M12 18h6"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(traces, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03) });

  const nodes = ['[cx="20"][cy="8"]', '[cx="18"]', '[cx="16"]', '[cx="20"][cy="21"]'].map((s) => el.querySelector(`circle${s}`));
  animate(nodes, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03, { startDelay: τmin }) });
}
