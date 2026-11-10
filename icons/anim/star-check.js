import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const star = el.querySelector('path[d^="m19.06"]');
  animate(star, { scale: [0.5, 1.15, 1], rotate: [0, 360] }, { duration: τmid });
  animate(star, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const tick = el.querySelector('path[d^="m15 18"]');
  animate(tick, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  animate(tick, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
