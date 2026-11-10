import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = ["M5 3a", "M9 3h", "M14 3h", "M19 3a", "M21 9v", "M21 14v", "M21 19a", "M14 21h", "M9 21h", "M5 21a", "M3 14v", "M3 9v"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmin, delay: stagger(0.03) });

  const plus = el.querySelectorAll('path[d="M8 12h8"], path[d="M12 8v8"]');
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
