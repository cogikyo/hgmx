import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const ghost = el.querySelector("g");
  animate(ghost, { translateY: [0, -3, 1, 0], rotate: [0, -6, 3, 0] }, { duration: τmax });
  const eyes = el.querySelectorAll('path[d$="v1"]');
  animate(eyes, { scaleY: [1, 0.1, 1] }, { duration: τmin, delay: τmid });
}
