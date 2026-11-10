import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const lines = ["M3 12h18", "M9 3v18", "M15 3v18"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(lines, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
