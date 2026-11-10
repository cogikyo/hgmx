import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const envelope = el.querySelector('path[d^="M22 10.5"]');
  animate(envelope, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const flap = el.querySelector('path[d^="m22 7"]');
  animate(flap, { scaleY: [-1, 1.15, 1] }, { duration: τmid });

  const mark = el.querySelector("g");
  animate(mark, { rotate: [0, -18, 14, -8, 0] }, { duration: τmid, delay: τmin });

  const dot = el.querySelector('path[d^="M20 22"]');
  animate(dot, { scale: [0, 1.25, 1] }, { duration: τmin, delay: τmid });
}
