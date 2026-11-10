import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const tray = el.querySelector("path");
  animate(tray, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const floor = el.querySelector("polyline");
  animate(floor, { translateY: [-6, 2, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
