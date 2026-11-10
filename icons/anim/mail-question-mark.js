import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const envelope = el.querySelector('path[d^="M22 10.5"]');
  animate(envelope, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const flap = el.querySelector('path[d^="m22 7"]');
  animate(flap, { scaleY: [-1, 1.15, 1] }, { duration: τmid });

  const mark = el.querySelectorAll('path[d^="M18 15"], path[d^="M20 22"]');
  animate(mark, { rotateY: [0, 360], scale: [1, 1.35, 1] }, { duration: τmid, delay: τmin });
}
