import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const body = el.querySelector('path[d^="M2 21"]');
  animate(body, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const star = el.querySelector('path[d^="M16.051"]');
  animate(star, { scale: [0, 1.25, 1], rotate: [0, 360] }, { duration: τmid, delay: τmin });

  const head = el.querySelector('path[d^="M8.927"]');
  animate(head, { translateY: [1, -2, 0] }, { duration: τmin, delay: τmid });
}
