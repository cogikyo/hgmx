import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const envelope = el.querySelector('path[d^="M22 12.5"]');
  animate(envelope, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const flap = el.querySelector('path[d^="m22 7"]');
  animate(flap, { scaleY: [-1, 1.15, 1] }, { duration: τmid });

  const lens = el.querySelector("g");
  animate(lens, { scale: [1.35, 1], translateX: [0, -4, 2, 0], translateY: [0, -3, -2, 0] }, { duration: τmax });
}
