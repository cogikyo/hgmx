import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const briefcase = el.querySelector("g");
  animate(briefcase, { translateY: [0, -3, 1, 0] }, { duration: τmax });

  const body = el.querySelector("g g");
  animate(body, { translateY: [0, 1.5, -0.5, 0] }, { duration: τmid, delay: τmin });

  const cross = el.querySelectorAll('path[d="M12 11v4"], path[d="M14 13h-4"]');
  animate(cross, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
