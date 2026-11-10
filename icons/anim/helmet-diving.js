import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const helmet = el.querySelectorAll(":scope > path, g.porthole");
  animate(helmet, { translateY: [0, -2, 0.5, 0] }, { duration: τmax });

  const porthole = el.querySelector("g.porthole");
  animate(porthole, { scale: [1, 1.15, 1] }, { duration: τmin, delay: τmid });
}
