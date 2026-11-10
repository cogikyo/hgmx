import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const bell = el.querySelector("g");
  animate(bell, { rotate: [0, -18, 14, -8, 0] }, { duration: τmax });
  const clapper = el.querySelector('path[d^="M10.268"]');
  animate(clapper, { translateX: [0, 2, -2, 0] }, { duration: τmid, delay: τmin });
  const tick = el.querySelector('path[d^="m15 8"]');
  animate(tick, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  animate(tick, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
