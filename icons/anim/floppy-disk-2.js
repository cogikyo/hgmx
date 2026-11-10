import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const disk = el.querySelector("path");
  animate(disk, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const shutter = el.querySelector("rect");
  animate(shutter, { translateX: [0, 2, -0.5, 0] }, { duration: τmid, delay: τmin });
}
