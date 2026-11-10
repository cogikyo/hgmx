import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const shell = el.querySelectorAll('ellipse[rx="10"], path[d^="M2 9"]');
  animate(shell, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const straw = el.querySelector('path[d^="m14 8"]');
  animate(straw, { translateX: [3, -0.5, 0], translateY: [-3, 0.5, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
