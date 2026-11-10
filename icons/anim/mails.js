import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const back = el.querySelector('path[d^="M17 19"]');
  animate(back, { translateX: [5, -1, 0], translateY: [-5, 1, 0], opacity: [0, 1] }, { duration: τmax });

  const envelope = el.querySelector("rect");
  animate(envelope, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const flap = el.querySelector('path[d^="m22 5.5"]');
  animate(flap, { scaleY: [-1, 1.15, 1] }, { duration: τmid, delay: τmin });
}
