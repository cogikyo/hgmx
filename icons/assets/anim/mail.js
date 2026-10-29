import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const envelope = el.querySelector("rect");
  animate(envelope, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const flap = el.querySelector("path");
  animate(flap, { scaleY: [-1, 1.15, 1] }, { duration: τmid, delay: τmin });
}
