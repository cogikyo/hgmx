import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = ["M5 3a", "M9 3h", "M14 3h", "M19 3a", "M21 9v", "M9 21h", "M5 21a", "M3 14v", "M3 9v"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmin, delay: stagger(0.03) });

  const bars = [el.querySelector('path[d="m21 16-5 5"]'), el.querySelector('path[d="m16 16 5 5"]')];
  animate(bars, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
