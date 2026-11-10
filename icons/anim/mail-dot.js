import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const envelope = el.querySelector('path[d^="M12.674"]');
  animate(envelope, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const flap = el.querySelector('path[d^="m2 7"]');
  animate(flap, { scaleY: [-1, 1.15, 1] }, { duration: τmid });

  const dot = el.querySelector("circle");
  animate(dot, { scale: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
