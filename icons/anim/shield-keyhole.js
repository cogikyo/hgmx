import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const shield = el.querySelector('path[d^="M20 13"]');
  animate(shield, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const keyhole = el.querySelector("g");
  animate(keyhole, { scale: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
