import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const sign = el.querySelector('path[d^="M2.586"]');
  animate(sign, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const minus = el.querySelector('path[d="M8 12h8"]');
  animate(minus, { scale: [0, 1.25, 1], rotate: [0, 360] }, { duration: τmid, delay: τmin });
}
