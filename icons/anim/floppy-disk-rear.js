import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const disk = el.querySelector('path[d^="M5 21"]');
  animate(disk, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const hub = el.querySelector("circle");
  animate(hub, { scale: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
