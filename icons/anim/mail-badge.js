import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const envelope = el.querySelector('path[d^="M22 7.7"]');
  animate(envelope, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const flap = el.querySelector('path[d^="M12 12.996"]');
  animate(flap, { scaleY: [-1, 1.15, 1] }, { duration: τmid });

  const emblem = el.querySelector("g");
  animate(emblem, { scale: [1, 0.85, 1.1, 1] }, { duration: τmid, delay: τmin });
}
