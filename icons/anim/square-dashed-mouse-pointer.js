import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = ["M5 3a", "M9 3h", "M14 3h", "M19 3a", "M21 9v", "M9 21h", "M5 21a", "M3 14v", "M3 9v"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmin, delay: stagger(0.03) });

  const pointer = el.querySelector('path[d^="M12.034"]');
  animate(pointer, { translateX: [0, -3, 1, 0], translateY: [0, -3, 1, 0] }, { duration: τmid, delay: τmin });
}
