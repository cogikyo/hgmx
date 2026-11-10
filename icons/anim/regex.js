import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const star = el.querySelectorAll('path[d^="M17 3"], path[d^="m12.67"]');
  animate(star, { rotate: [0, 180], scale: [1, 1.25, 1] }, { duration: τmid });

  const dot = el.querySelector('path[d^="M9 17"]');
  animate(dot, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
