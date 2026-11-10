import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const halves = [el.querySelector('path[d^="M7.298 3"]'), el.querySelector('path[d^="M7.298 20"]')];
  animate(halves, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const dashes = ["M4 12H2", "M10 12H8", "M16 12h-2", "M22 12h-2"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(dashes, { scaleX: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });
}
