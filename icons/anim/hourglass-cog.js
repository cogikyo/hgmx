import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const gear = el.querySelector("g.gear");
  animate(gear, { scale: [0, 1.25, 1], rotate: [0, 360] }, { duration: τmid, delay: τmin });
  const caps = el.querySelectorAll('path[d="M5 2h14"], path[d="M5 22h6.159"]');
  animate(caps, { scaleX: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
