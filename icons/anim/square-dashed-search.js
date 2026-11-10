import { animate, stagger, τmax, τmin } from "../motion.js";

export default function (el) {
  const frame = ["M5 3a", "M9 3h", "M14 3h", "M19 3a", "M21 9v", "M9 21h", "M5 21a", "M3 14v", "M3 9v"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmin, delay: stagger(0.03) });

  const lens = el.querySelector("g");
  animate(lens, { scale: [1.35, 1] }, { duration: τmax });
  animate(lens, { translateX: [0, -3, 2, 0], translateY: [0, -3, -1, 0] }, { duration: τmax });
}
