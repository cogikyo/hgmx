import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const front = el.querySelector("rect");
  animate(front, { scale: [1, 1.3, 1] }, { duration: τmax });

  const back = el.querySelector('path[d^="M4 10"]');
  animate(back, { translateX: [0, -1.5, 0], translateY: [0, -1.5, 0] }, { duration: τmid });

  const behind = el.querySelector('path[d^="M14 20"]');
  animate(behind, { translateX: [0, 1.5, 0], translateY: [0, 1.5, 0] }, { duration: τmid });
}
