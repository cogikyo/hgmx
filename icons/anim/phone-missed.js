import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const handset = el.querySelector('path[d^="M13.832"]');
  animate(handset, { rotate: [0, -10, 6, 0] }, { duration: τmid });

  const x = el.querySelectorAll('path[d^="m16 2"], path[d^="m22 2"]');
  animate(x, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
