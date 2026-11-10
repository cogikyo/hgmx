import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const pin = el.querySelector('path[d^="M20 10c"]');
  animate(pin, { translateY: [-6, 1, 0], scaleY: [1.15, 0.9, 1] }, { duration: τmid });

  const tick = el.querySelector('path[d^="m9 10 "]');
  animate(tick, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  animate(tick, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
