import { animate, stagger, τmin } from "../motion.js";

export default function (el) {
  const frame = ["M4 3a", "M8 3h", "M14 3h", "M22 6V", "M22 10v", "M20 19a", "M16 19h", "M8 19h", "M2 16v", "M2 12v"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmin, delay: stagger(0.03) });
}
