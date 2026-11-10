import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const envelope = el.querySelector('path[d^="M6.835"]');
  animate(envelope, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const flap = el.querySelector("g");
  animate(flap, { scaleY: [-1, 1.15, 1] }, { duration: τmid });

  const hands = el.querySelector('path[d^="M16 14"]');
  animate(hands, { rotate: [0, 360] }, { duration: τmid, delay: τmin });
}
