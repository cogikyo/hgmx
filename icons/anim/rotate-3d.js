import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const meridian = el.querySelector('path[d^="M16.47214"]');
  animate(meridian, { rotateY: [0, 360] }, { duration: τmax });

  const head = el.querySelector('path[d^="m15.194"]');
  animate(head, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
