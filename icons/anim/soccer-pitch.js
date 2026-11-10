import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const spot = el.querySelector("circle");
  animate(spot, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3 });
  const halfway = [el.querySelector('path[d="M12 5v5"]'), el.querySelector('path[d="M12 14v5"]')];
  animate(halfway, { scaleY: [0, 1] }, { duration: τmid, delay: τmin });
}
