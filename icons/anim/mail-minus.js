import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const envelope = el.querySelector('path[d^="M22 15"]');
  animate(envelope, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const flap = el.querySelector('path[d^="m22 7"]');
  animate(flap, { scaleY: [-1, 1.15, 1] }, { duration: τmid });

  const minus = el.querySelector('path[d^="M16 19"]');
  animate(minus, { scale: [0, 1.25, 1], rotate: [0, 360] }, { duration: τmid, delay: τmin });
}
