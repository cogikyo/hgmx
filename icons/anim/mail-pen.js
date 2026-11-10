import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const envelope = el.querySelector('path[d^="M22 10.38"]');
  animate(envelope, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const flap = el.querySelector('path[d^="m22 7"]');
  animate(flap, { scaleY: [-1, 1.15, 1] }, { duration: τmid });

  const pen = el.querySelector('path[d^="M15.363"]');
  animate(pen, { rotate: [-15, 15, -10, 0], translateX: [0, -3, 0] }, { duration: τmax });
}
