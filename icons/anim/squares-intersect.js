import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const first = el.querySelectorAll('path[d^="M14 2"], path[d^="M2 10"], path[d^="M2 4"], path[d^="M4 16"], path[d^="M8 2"]');
  animate(first, { translateX: [-2, 0.5, 0], translateY: [-2, 0.5, 0] }, { duration: τmid });

  const second = el.querySelectorAll('path[d^="M10 22"], path[d^="M16 22"], path[d^="M20 8"], path[d^="M22 14"], path[d^="M22 20"]');
  animate(second, { translateX: [2, -0.5, 0], translateY: [2, -0.5, 0] }, { duration: τmid });

  const overlap = el.querySelector('path[d^="M8 10"]');
  animate(overlap, { scale: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
