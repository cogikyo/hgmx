import { animate, τmax } from "../motion.js";

export default function (el) {
  const pupil = el.querySelector('path[d="M14 7h.01"]');
  animate(pupil, { translateX: [0, 0.5, 1.5, 0.5, 0], translateY: [0, -1, 0, 1, 0] }, { duration: τmax });
}
