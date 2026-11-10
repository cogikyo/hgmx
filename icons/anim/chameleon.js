import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const eye = el.querySelector("circle");
  animate(eye, { scale: [1, 1.25, 1] }, { duration: τmid });
  const pupil = el.querySelector('path[d="M14 7h.01"]');
  animate(pupil, { translateX: [0, 2, -1, 2, 0], translateY: [0, -1.5, 1.5, 1, 0] }, { duration: τmax });
}
