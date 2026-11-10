import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const ghost = el.querySelector("g");
  animate(ghost, { translateX: [0, 3, -1, 0] }, { duration: τmax });
  const eyes = el.querySelectorAll('path[d$="h.01"]');
  animate(eyes, { translateX: [0, 1.5, 0] }, { duration: τmid });
}
